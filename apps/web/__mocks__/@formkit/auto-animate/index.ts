const autoAnimate = jest.fn((el: Element) => ({
  parent: el,
  enable: jest.fn(),
  disable: jest.fn(),
  isEnabled: jest.fn(),
}));

export default autoAnimate;
