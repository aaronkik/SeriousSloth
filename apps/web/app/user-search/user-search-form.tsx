'use client';

import { useSearchParams } from 'next/navigation';
import { Suspense, useActionState } from 'react';
import { useFormStatus } from 'react-dom';
import {
  type SearchFormState,
  searchUsernameAction,
} from '~/app/user-search/actions';
import { usernameSchema } from '~/app/user-search/schemas';
import { Button } from '~/components/ui/button';
import { Field, FieldDescription } from '~/components/ui/field';
import { Input } from '~/components/ui/input';
import { Spinner } from '~/components/ui/spinner';

const initialState: SearchFormState = {};

const SearchButton = () => {
  const { pending } = useFormStatus();

  return (
    <Button className='min-w-20' disabled={pending} type='submit'>
      {pending ? <Spinner className='size-4' /> : 'Search'}
    </Button>
  );
};

interface FormProps {
  defaultUsername?: string;
}

const UserSearchFormUI = ({ defaultUsername }: FormProps) => {
  const [state, formAction] = useActionState(
    searchUsernameAction,
    initialState,
  );
  const hasFormErrors = !!state.errors?.length;

  return (
    <form action={formAction} role='search'>
      <Field orientation='horizontal' data-invalid={hasFormErrors}>
        <Input
          aria-invalid={hasFormErrors}
          autoComplete='off'
          defaultValue={defaultUsername}
          minLength={usernameSchema.minLength ?? 0}
          maxLength={usernameSchema.maxLength ?? 0}
          name='username'
          placeholder='Search...'
          required
          type='search'
        />
        <SearchButton />
      </Field>
      {state.errors ? (
        <div className='pt-2 px-2'>
          {state.errors.map((error) => (
            <FieldDescription className='first:mt-2' key={error}>
              {error}
            </FieldDescription>
          ))}
        </div>
      ) : null}
    </form>
  );
};

const UserSearchFormWithParams = () => {
  const searchParams = useSearchParams();
  const defaultUsername = searchParams?.get('username') ?? undefined;
  return <UserSearchFormUI defaultUsername={defaultUsername} />;
};

const UserSearchForm = () => (
  <Suspense fallback={<UserSearchFormUI />}>
    <UserSearchFormWithParams />
  </Suspense>
);

export default UserSearchForm;
