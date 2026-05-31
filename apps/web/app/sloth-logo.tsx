import Image from 'next/image';
import { ComponentPropsWithoutRef } from 'react';
import slothFace from './sloth-face-square.png';

type Props = Omit<
  Partial<ComponentPropsWithoutRef<typeof Image>>,
  'alt' | 'src' | 'placeholder'
>;

const SlothLogo = (props: Props) => (
  <Image {...props} alt='Sloth logo' src={slothFace} placeholder='blur' />
);

export default SlothLogo;
