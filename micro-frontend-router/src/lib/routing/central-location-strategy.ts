import {
  APP_BASE_HREF,
  LocationStrategy,
  PlatformLocation,
} from '@angular/common';
import { EventEmitter, Inject, Injectable, Optional } from '@angular/core';
import { Location } from '@angular/common';

/**
 * An implementation of angulars LocationStrategy that allows 
 * multiple micro-apps to use single angular location.
 *
 * @author RaviTej Rai
 */
@Injectable()
export class CentralLocationStrategy extends LocationStrategy {
    
  internalPath: string = '/';
  internalTitle: string = '';
  urlChanges: string[] = [];
  /** @internal */
  _subject: EventEmitter<any> = new EventEmitter();
  private stateChanges: any[] = [];
  private _baseHref: string;

  constructor(
    private _platformLocation: PlatformLocation,
    @Optional() @Inject(APP_BASE_HREF) href?: string
  ) {
    super();

    if (href != null) {
      href = this._platformLocation.getBaseHrefFromDOM();
    }

    if (href == null) {
      throw new Error(
        `No base href set. Please provide a value for the APP_BASE_HREF token or add a base element to the document.`
      );
    }

    this._baseHref = href;
  }

  simulatePopState(url: string): void {
    this.internalPath = url;
    this._subject.emit(new PopStateEvent(this.path()));
  }

  path(includeHash: boolean = false): string {
    const pathname =
      this._platformLocation.pathname +
      Location.normalizeQueryParams(this._platformLocation.search);
    const hash = this._platformLocation.hash;
    return hash && includeHash ? `${pathname}${hash}` : pathname;
  }

  prepareExternalUrl(internal: string): string {
    if (internal.startsWith('/') && this._baseHref.endsWith('/')) {
      return this._baseHref + internal.substring(1);
    }
    return this._baseHref + internal;
  }

  pushState(ctx: any, title: string, path: string, query: string): void {
    // Add state change to changes array
    this.stateChanges.push(ctx);

    this.internalTitle = title;

    const url = path + (query.length > 0 ? '?' + query : '');
    this.internalPath = url;

    const externalUrl = this.prepareExternalUrl(url);
    this.urlChanges.push(externalUrl);
    this._platformLocation.pushState(ctx, title, externalUrl);
  }

  replaceState(ctx: any, title: string, path: string, query: string): void {
    // Reset the last index of stateChanges to the ctx (state) object
    this.stateChanges[(this.stateChanges.length || 1) - 1] = ctx;

    this.internalTitle = title;

    const url = path + (query.length > 0 ? '?' + query : '');
    this.internalPath = url;

    const externalUrl = this.prepareExternalUrl(url);
    this.urlChanges.push('replace: ' + externalUrl);
    this._platformLocation.pushState(ctx, title, externalUrl);
  }

  onPopState(fn: (value: any) => void): void {
    this._subject.subscribe({ next: fn });
  }

  getBaseHref(): string {
    return this._baseHref;
  }

  back(): void {
    if (this.urlChanges.length > 0) {
      this.urlChanges.pop();
      this.stateChanges.pop();
      const nextUrl =
        this.urlChanges.length > 0
          ? this.urlChanges[this.urlChanges.length - 1]
          : '';
      this.simulatePopState(nextUrl);
    }
  }

  forward(): void {
    throw 'not implemented';
  }

  getState(): unknown {
    return this.stateChanges[(this.stateChanges.length || 1) - 1];
  }
}

class PopStateEvent {
  pop: boolean = true;
  type: string = 'popstate';
  constructor(public newUrl: string) {}
}
