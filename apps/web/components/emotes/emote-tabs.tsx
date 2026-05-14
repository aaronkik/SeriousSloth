import { useState } from 'react';
import type { ActiveEmote, RemovedEmote } from '~/lib/api/emotes-service';
import EmotesList from './emotes-list';

type Tab = 'active' | 'removed';

type Props = {
  activeEmotes: ActiveEmote[];
  removedEmotes: RemovedEmote[];
};

const tabBase =
  'px-4 py-2 rounded-full text-sm font-medium transition-colors flex items-center gap-2';
const tabIdle = 'text-gray-600 hover:bg-gray-100';
const tabActive = 'bg-purple-600 text-white';
const countBase = 'rounded-full px-2 py-0.5 text-xs font-semibold';
const countIdle = 'bg-gray-200 text-gray-700';
const countActive = 'bg-white/20 text-white';

const EmoteTabs = ({ activeEmotes, removedEmotes }: Props) => {
  const [tab, setTab] = useState<Tab>('active');

  const tabs: Array<{ id: Tab; label: string; count: number }> = [
    { id: 'active', label: 'Active', count: activeEmotes.length },
    { id: 'removed', label: 'Removed', count: removedEmotes.length },
  ];

  const emotes = tab === 'active' ? activeEmotes : removedEmotes;

  return (
    <div>
      <div
        role='tablist'
        className='mt-6 flex justify-center gap-2'
        data-testid='emoteTabs'
      >
        {tabs.map(({ id, label, count }) => {
          const selected = tab === id;
          return (
            <button
              key={id}
              role='tab'
              type='button'
              aria-selected={selected}
              onClick={() => setTab(id)}
              className={`${tabBase} ${selected ? tabActive : tabIdle}`}
              data-testid={`emoteTab-${id}`}
            >
              <span>{label}</span>
              <span
                className={`${countBase} ${selected ? countActive : countIdle}`}
                data-testid={`emoteTabCount-${id}`}
              >
                {count}
              </span>
            </button>
          );
        })}
      </div>
      <EmotesList emotes={emotes} />
    </div>
  );
};

export default EmoteTabs;
