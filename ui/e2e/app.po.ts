export class SmokePage {
  navigateTo() {
    return browser.get('/');
  }

  getParagraphText() {
    return element(by.css('smoke-app h1')).getText();
  }
}
