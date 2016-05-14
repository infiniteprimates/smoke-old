import {
  beforeEachProviders,
  describe,
  expect,
  it,
  inject
} from '@angular/core/testing';
import { SmokeAppComponent } from '../app/smoke.component';

beforeEachProviders(() => [SmokeAppComponent]);

describe('App: Smoke', () => {
  it('should create the app',
      inject([SmokeAppComponent], (app: SmokeAppComponent) => {
    expect(app).toBeTruthy();
  }));

  it('should have as title \'smoke works!\'',
      inject([SmokeAppComponent], (app: SmokeAppComponent) => {
    expect(app.title).toEqual('smoke works!');
  }));
});
