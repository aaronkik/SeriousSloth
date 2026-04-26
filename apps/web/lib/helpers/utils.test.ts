import { capitaliseFirstLetter } from './utils';

describe('capitaliseFirstLetter', () => {
  it('Returns the capitalised first letter of each string', () => {
    const testValues = [
      {
        initial: 'abcdefg',
        expected: 'Abcdefg',
      },
      {
        initial: 'A test string',
        expected: 'A test string',
      },
      {
        initial: 'a bcd efg',
        expected: 'A bcd efg',
      },
      {
        initial: 'a',
        expected: 'A',
      },
    ];

    for (const testValue of testValues) {
      const { initial, expected } = testValue;
      expect(capitaliseFirstLetter(initial)).toBe(expected);
    }
  });
});
