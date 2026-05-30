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
    <Button className='min-w-20' disabled={pending} type='submit'>
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
      <Field orientation='horizontal' data-invalid={hasFormErrors}>
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

export default UserSearchForm;
