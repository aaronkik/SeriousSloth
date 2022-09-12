import Image from 'next/future/image';
import { ComponentPropsWithoutRef } from 'react';
import slothFace from '~/public/assets/sloth-face-square.png';

type Props = Omit<
  Partial<ComponentPropsWithoutRef<typeof Image>>,
  'alt' | 'src' | 'placeholder'
>;

const SlothLogo = (props: Props) => (
  <Image {...props} alt='Sloth logo' src={slothFace} placeholder='blur' />
);

export default SlothLogo;
