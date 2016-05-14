import { SmokePage } from './app.po';

describe('smoke App', function() {
  let page: SmokePage;

  beforeEach(() => {
    page = new SmokePage();
  })

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('smoke works!');
  });
});
