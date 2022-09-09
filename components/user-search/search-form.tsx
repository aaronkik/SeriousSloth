import autoAnimate from '@formkit/auto-animate';
import { useEffect, useRef } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { Button, FormErrorMessage, Input, Spinner } from '~/components/shared';
import { usernameFormRules } from '~/constants/form';
import { fetchUsers } from '~/lib/api';
import { GetUsers } from '~/types/twitch';

interface Props {
  setUserResponse: React.Dispatch<React.SetStateAction<GetUsers | undefined>>;
}

const SearchForm = ({ setUserResponse }: Props) => {
  const formRef = useRef(null);

  useEffect(() => {
    formRef.current && autoAnimate(formRef.current);
  }, []);

  const {
    formState: { errors, isSubmitting },
    handleSubmit,
    register,
  } = useForm<{ username: string }>({
    defaultValues: {
      username: '',
    },
  });

  const submitUsername = handleSubmit(async ({ username }) => {
    try {
      const response = await fetchUsers({ username });
      setUserResponse(response);
    } catch (error: any) {
      console.error(error);
      toast(
        <div data-testid='userError'>{error?.message || 'Unknown error'}</div>
      );
    }
  });

  return (
    <form
      id='user-search-form'
      onSubmit={submitUsername}
      ref={formRef}
      role='search'
    >
      <div className='flex gap-4'>
        <Input
          aria-label='user-search'
          placeholder='Twitch'
          type='search'
          {...register('username', usernameFormRules)}
        />
        <Button aria-label='search' className='min-w-[5rem]' type='submit'>
          {isSubmitting ? <Spinner className='h-4 w-4' /> : 'Search'}
        </Button>
      </div>
      {errors.username && (
        <FormErrorMessage className='mt-2' data-testid='searchErrorMessage'>
          {errors.username.message}
        </FormErrorMessage>
      )}
    </form>
  );
};

export default SearchForm;
