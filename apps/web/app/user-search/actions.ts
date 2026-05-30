'use server';

import { redirect } from 'next/navigation';
import { usernameSchema } from '~/app/user-search/schemas';

export interface SearchFormState {
  errors?: string[];
}

export const searchUsernameAction = async (
  _prev: SearchFormState,
  formData: FormData,
): Promise<SearchFormState> => {
  const username = formData.get('username');
  const parsed = usernameSchema.safeParse(username);

  if (!parsed.success) {
    return { errors: parsed.error.issues.map((issue) => issue.message) };
  }

  // Page server-fetches from searchParams
  redirect(`/user-search?username=${encodeURIComponent(parsed.data)}`);
};
