import { useForm } from 'react-hook-form';
import { usernameFormRules } from '~/constants/form';
import { fetchUsers } from '~/lib/api';
import { toast } from 'react-toastify';
import { GetUsers } from '~/types/twitch';

interface Props {
  setUserResponse: React.Dispatch<React.SetStateAction<GetUsers | undefined>>;
}

const SearchForm = ({ setUserResponse }: Props) => {
  const {
    formState: { errors },
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
    <form id='user-search-form' role='search' onSubmit={submitUsername}>
      <input
        aria-label='user-search'
        type='search'
        {...register('username', usernameFormRules)}
      />
      <button aria-label='search' type='submit'>
        Search
      </button>
      {errors.username && (
        <p data-testid='searchErrorMessage' role='alert'>
          {errors.username.message}
        </p>
      )}
    </form>
  );
};

export default SearchForm;
