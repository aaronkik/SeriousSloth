'use client';

import { useActionState } from 'react';
import { useFormStatus } from 'react-dom';
import {
  type SearchFormState,
  searchUsernameAction,
} from '~/app/user-search/actions';
import { Button } from '~/components/ui/button';
import { Field, FieldDescription } from '~/components/ui/field';
import { Input } from '~/components/ui/input';
import { Spinner } from '~/components/ui/spinner';

const initialState: SearchFormState = {};

const SearchButton = () => {
  const { pending } = useFormStatus();

  return (
    <Button
      className='min-w-20'
      aria-disabled={pending}
      type={pending ? 'button' : 'submit'}
    >
      {pending ? <Spinner className='size-4' /> : 'Search'}
    </Button>
  );
};

interface Props {
  defaultUsername?: string;
}

const UserSearchForm = ({ defaultUsername }: Props) => {
  const [state, formAction] = useActionState(
    searchUsernameAction,
    initialState,
  );
  const hasFormErrors = !!state.errors?.length;

  return (
    <form action={formAction} role='search'>
      <Field
        className='pb-2'
        orientation='horizontal'
        data-invalid={hasFormErrors}
      >
        <Input
          aria-invalid={hasFormErrors}
          autoComplete='off'
          defaultValue={defaultUsername}
          minLength={3}
          maxLength={120}
          name='username'
          placeholder='Search...'
          required
          type='search'
        />
        <SearchButton />
      </Field>
      {state.errors?.map((error) => (
        <FieldDescription key={error}>{error}</FieldDescription>
      ))}
    </form>
  );
};

export default UserSearchForm;
