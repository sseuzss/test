import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { Store } from '@ngrx/store';
import { catchError, forkJoin, map, of, switchMap, withLatestFrom } from 'rxjs';

import {
    AgentsService,
    allGroupedListQuery,
    defaultListQuery,
    ErrorResponse,
    EventsService,
    EventsSQLMappers,
    GroupedData,
    GroupsService,
    ModelsModuleA,
    ModulesService,
    PoliciesService,
    PrivateEvents,
    PrivatePolicies,
    PrivateSystemShortModules,
    SuccessResponse
} from '@soldr/api';
import { ViewMode } from '@soldr/shared';
import { PoliciesFacade } from '@soldr/store/policies';

import * as ModulesInstancesActions from './modules-instances.actions';
import { State } from './modules-instances.reducer';
import {
    selectEntityId,
    selectEventsGridFiltration,
    selectEventsGridSearch,
    selectEventsGridSorting,
    selectModuleName,
    selectViewMode
} from './modules-instances.selectors';

@Injectable()
export class ModulesInstancesEffects {
    fetchModule$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.fetchModule),
            withLatestFrom(this.store.select(selectViewMode), this.store.select(selectModuleName)),
            switchMap(([{ entityHash }, viewMode, moduleName]) =>
                this.getEntityService(viewMode)
                    .fetchModule(entityHash, moduleName)
                    .pipe(
                        map((response: SuccessResponse<ModelsModuleA>) =>
                            ModulesInstancesActions.fetchModuleSuccess({ module: response.data })
                        ),
                        catchError(() => of(ModulesInstancesActions.fetchModuleFailure()))
                    )
            )
        )
    );

    fetchModulePolicy$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.fetchModuleSuccess),
            withLatestFrom(this.store.select(selectViewMode)),
            switchMap(([{ module }]) =>
                this.policiesService
                    .fetchList(
                        defaultListQuery({
                            pageSize: 1,
                            filters: [
                                {
                                    field: 'id',
                                    value: module.policy_id
                                }
                            ]
                        })
                    )
                    .pipe(
                        map((response: SuccessResponse<PrivatePolicies>) =>
                            ModulesInstancesActions.fetchPolicySuccess({ data: response.data })
                        ),
                        catchError(() => of(ModulesInstancesActions.fetchPolicyFailure()))
                    )
            )
        )
    );

    fetchEvents$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.fetchEvents),
            withLatestFrom(
                this.store.select(selectViewMode),
                this.store.select(selectEntityId),
                this.store.select(selectModuleName),
                this.store.select(selectEventsGridSearch),
                this.store.select(selectEventsGridFiltration),
                this.store.select(selectEventsGridSorting)
            ),
            switchMap(([action, viewMode, entityId, moduleName, search, filtration, sorting]) => {
                const page = action.page || 1;
                const baseFiltration = this.getEventsFiltration(viewMode, entityId, moduleName);

                return this.eventsService
                    .fetchEvents(
                        defaultListQuery({
                            page,
                            filters: [
                                ...baseFiltration,
                                {
                                    field: 'data',
                                    value: search
                                },
                                ...filtration
                            ],
                            sort: sorting || {}
                        })
                    )
                    .pipe(
                        map((eventsResponse) =>
                            ModulesInstancesActions.fetchEventsSuccess({
                                data: (eventsResponse as SuccessResponse<PrivateEvents>).data,
                                page
                            })
                        ),
                        catchError(() => of(ModulesInstancesActions.fetchEventsFailure()))
                    );
            })
        )
    );

    fetchModuleEventsFilterItems$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.fetchModuleEventsFilterItems),
            withLatestFrom(
                this.store.select(selectViewMode),
                this.store.select(selectEntityId),
                this.store.select(selectModuleName),
                this.store.select(selectEventsGridFiltration)
            ),
            switchMap(([_, viewMode, entityId, moduleName, filtration]) =>
                forkJoin([
                    this.eventsService.fetchEvents(
                        allGroupedListQuery(
                            {
                                filters: [...this.getEventsFiltration(viewMode, entityId, moduleName), ...filtration]
                            },
                            EventsSQLMappers.AgentId
                        )
                    ),
                    this.eventsService.fetchEvents(
                        allGroupedListQuery(
                            {
                                filters: [...this.getEventsFiltration(viewMode, entityId, moduleName), ...filtration]
                            },
                            EventsSQLMappers.GroupId
                        )
                    )
                ]).pipe(
                    map(
                        ([agentsResponse, policiesResponse]: [
                            SuccessResponse<GroupedData>,
                            SuccessResponse<GroupedData>
                        ]) =>
                            ModulesInstancesActions.fetchModuleEventsFilterItemsSuccess({
                                agentIds: agentsResponse.data.grouped,
                                groupIds: policiesResponse.data.grouped
                            })
                    ),
                    catchError((error: ErrorResponse) =>
                        of(ModulesInstancesActions.fetchModuleEventsFilterItemsFailure({ error }))
                    )
                )
            )
        )
    );

    setEventsGridFiltration$ = createEffect(() =>
        this.actions$.pipe(
            ofType(
                ...[
                    ModulesInstancesActions.setEventsGridSearch,
                    ModulesInstancesActions.setEventsGridFiltration,
                    ModulesInstancesActions.resetEventsFiltration,
                    ModulesInstancesActions.setEventsGridSorting
                ]
            ),
            switchMap(() => [
                ModulesInstancesActions.fetchEvents({ page: 1 }),
                ModulesInstancesActions.fetchModuleEventsFilterItems()
            ])
        )
    );

    fetchModuleVersions$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.fetchModuleVersions),
            switchMap(({ moduleName }) =>
                this.modulesService.fetchVersions(moduleName).pipe(
                    map((response: SuccessResponse<PrivateSystemShortModules>) =>
                        ModulesInstancesActions.fetchModuleVersionSuccess({ data: response.data })
                    ),
                    catchError(() => of(ModulesInstancesActions.fetchModuleVersionsFailed()))
                )
            )
        )
    );

    enableModule$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.enableModule),
            switchMap(({ policyHash, moduleName }) =>
                this.policiesService.activateModule(policyHash, moduleName).pipe(
                    map(() => ModulesInstancesActions.enableModuleSuccess()),
                    catchError(({ error }) => of(ModulesInstancesActions.enableModuleFailure({ error })))
                )
            )
        )
    );

    disableModule$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.disableModule),
            switchMap(({ policyHash, moduleName }) =>
                this.policiesService.deactivateModule(policyHash, moduleName).pipe(
                    map(() => ModulesInstancesActions.disableModuleSuccess()),
                    catchError(({ error }) => of(ModulesInstancesActions.disableModuleFailure({ error })))
                )
            )
        )
    );

    deleteModule$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.deleteModule),
            switchMap(({ policyHash, moduleName }) =>
                this.policiesService.deleteModule(policyHash, moduleName).pipe(
                    map(() => ModulesInstancesActions.deleteModuleSuccess()),
                    catchError(({ error }) => of(ModulesInstancesActions.deleteModuleFailure({ error })))
                )
            )
        )
    );

    changeModuleVersion$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.changeModuleVersion),
            switchMap(({ policyHash, moduleName, version }) =>
                this.policiesService.changeModuleVersion(policyHash, moduleName, version).pipe(
                    map(() => ModulesInstancesActions.changeModuleVersionSuccess()),
                    catchError(({ error }) => of(ModulesInstancesActions.changeModuleVersionFailure({ error })))
                )
            )
        )
    );

    updateModule$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.updateModule),
            switchMap(({ policyHash, moduleName, version }) =>
                this.policiesService.updateModule(policyHash, moduleName, version).pipe(
                    map(() => ModulesInstancesActions.updateModuleSuccess()),
                    catchError(({ error }) => of(ModulesInstancesActions.updateModuleFailure({ error })))
                )
            )
        )
    );

    saveModule$ = createEffect(() =>
        this.actions$.pipe(
            ofType(ModulesInstancesActions.saveModuleConfig),
            switchMap(({ module, policyHash }) =>
                this.policiesService.storeModule(policyHash, module).pipe(
                    map(() => ModulesInstancesActions.saveModuleConfigSuccess()),
                    catchError(({ error }) => of(ModulesInstancesActions.saveModuleConfigFailure({ error })))
                )
            )
        )
    );

    constructor(
        private agentsService: AgentsService,
        private groupsService: GroupsService,
        private policiesService: PoliciesService,
        private eventsService: EventsService,
        private modulesService: ModulesService,
        private policiesFacade: PoliciesFacade,
        private actions$: Actions,
        private store: Store<State>
    ) {}

    private getEntityService(viewMode: ViewMode) {
        switch (viewMode) {
            case ViewMode.Agents:
                return this.agentsService;
            case ViewMode.Groups:
                return this.groupsService;
            case ViewMode.Policies:
                return this.policiesService;
        }
    }

    private getEventsFiltration(viewMode: ViewMode, entityId: number, moduleName: string) {
        switch (viewMode) {
            case ViewMode.Agents:
                return [
                    { field: 'agent_id', value: entityId },
                    { field: 'module_name', value: moduleName }
                ];
            case ViewMode.Groups:
                return [
                    { field: 'group_id', value: entityId },
                    { field: 'module_name', value: moduleName }
                ];
            case ViewMode.Policies:
                return [
                    { field: 'policy_id', value: entityId },
                    { field: 'module_name', value: moduleName }
                ];
        }
    }
}
