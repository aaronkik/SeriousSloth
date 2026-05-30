'use client';

import { useState } from 'react';
import { Badge } from '~/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs';
import type {
  ActiveEmoteEntry,
  RemovedEmoteEntry,
} from '~/lib/api/emotes-service';
import EmotesList from './emotes-list';

type Tab = 'active' | 'removed';

type Props = {
  activeEmotes: ActiveEmoteEntry[];
  removedEmotes: RemovedEmoteEntry[];
};

const EmoteTabs = ({ activeEmotes, removedEmotes }: Props) => {
  const [tab, setTab] = useState<Tab>('active');

  const tabs: Array<{ id: Tab; label: string; count: number; emotes: (ActiveEmoteEntry | RemovedEmoteEntry)[]; emptyMessage: string }> = [
    { id: 'active', label: 'Active', count: activeEmotes.length, emotes: activeEmotes, emptyMessage: 'No active emotes' },
    { id: 'removed', label: 'Removed', count: removedEmotes.length, emotes: removedEmotes, emptyMessage: 'No removed emotes' },
  ];

  return (
    <Tabs
      value={tab}
      onValueChange={(v) => setTab(v as Tab)}
      className='mt-6'
      data-testid='emoteTabs'
    >
      <TabsList className='mx-auto'>
        {tabs.map(({ id, label, count }) => (
          <TabsTrigger key={id} value={id} data-testid={`emoteTab-${id}`}>
            {label}
            <Badge variant='secondary' data-testid={`emoteTabCount-${id}`}>
              {count}
            </Badge>
          </TabsTrigger>
        ))}
      </TabsList>
      {tabs.map(({ id, emotes, emptyMessage }) => (
        <TabsContent key={id} value={id}>
          <EmotesList emotes={emotes} emptyMessage={emptyMessage} />
        </TabsContent>
      ))}
    </Tabs>
  );
};

export default EmoteTabs;
