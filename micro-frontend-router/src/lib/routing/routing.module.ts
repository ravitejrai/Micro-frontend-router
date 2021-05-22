import { CentralLocationStrategy } from './central-location-strategy';
import {
  NgModule,
  NgModuleFactoryLoader,
  Compiler,
  Injector,
  Optional,
  ModuleWithProviders,
  SystemJsNgModuleLoader,
} from '@angular/core';
import {
  RouterModule,
  Router,
  UrlSerializer,
  ChildrenOutletContexts,
  ROUTES,
  ROUTER_CONFIGURATION,
  UrlHandlingStrategy,
  PreloadingStrategy,
  NoPreloading,
  provideRoutes,
  Routes,
  ExtraOptions,
  Route,
} from '@angular/router';
import { ɵROUTER_PROVIDERS as ROUTER_PROVIDERS } from '@angular/router';
import { LocationStrategy, Location } from '@angular/common';
import { flatten } from '@angular/compiler';

export function setupMicroRouter(
  urlSerializer: UrlSerializer,
  contexts: ChildrenOutletContexts,
  location: Location,
  loader: NgModuleFactoryLoader,
  compiler: Compiler,
  injector: Injector,
  routes: Route[][],
  opts?: ExtraOptions,
  urlHandlingStrategy?: UrlHandlingStrategy
): Router {
  // tslint:disable-next-line:no-non-null-assertion
  return new Router(
    null!,
    urlSerializer,
    contexts,
    location,
    injector,
    loader,
    compiler,
    flatten(routes)
  );
}

@NgModule({
  exports: [RouterModule],
  providers: [
    ROUTER_PROVIDERS,
    { provide: Location, useClass: Location },
    { provide: LocationStrategy, useClass: CentralLocationStrategy },
    { provide: NgModuleFactoryLoader, useClass: SystemJsNgModuleLoader },
    {
      provide: Router,
      useFactory: setupMicroRouter,
      deps: [
        UrlSerializer,
        ChildrenOutletContexts,
        Location,
        NgModuleFactoryLoader,
        Compiler,
        Injector,
        ROUTES,
        ROUTER_CONFIGURATION,
        [UrlHandlingStrategy, new Optional()],
      ],
    },
    { provide: PreloadingStrategy, useExisting: NoPreloading },
    provideRoutes([]),
  ],
})

export class RoutingModule {
  static withRoutes(
    routes: Routes,
    config?: ExtraOptions
  ): ModuleWithProviders<RoutingModule> {
    return {
      ngModule: RoutingModule,
      providers: [
        provideRoutes(routes),
        { provide: ROUTER_CONFIGURATION, useValue: config ? config : {} },
      ],
    };
  }
}
