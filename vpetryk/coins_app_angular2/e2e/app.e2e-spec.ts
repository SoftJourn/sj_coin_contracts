import { CoinsPage } from './app.po';

describe('coins App', () => {
  let page: CoinsPage;

  beforeEach(() => {
    page = new CoinsPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});
