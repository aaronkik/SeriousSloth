'use client';

import { type ReactNode, useState } from 'react';
import { Badge } from '~/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs';
import type { EmoteListEntry, EmotesByDate } from '~/lib/api/emotes-service';
import { countEmotes } from '~/lib/helpers';
import EmotesList, { EmoteDateGroups, EmoteGrid } from './emotes-list';

type Tab = 'active' | 'added' | 'removed';

type Props = {
  activeEmotes: EmoteListEntry[];
  addedEmotes: EmotesByDate;
  removedEmotes: EmotesByDate;
};

const EmoteTabs = ({ activeEmotes, addedEmotes, removedEmotes }: Props) => {
  const [tab, setTab] = useState<Tab>('active');

  const tabs: Array<{
    id: Tab;
    label: string;
    count: number;
    emptyMessage: string;
    content: ReactNode;
  }> = [
    {
      id: 'active',
      label: 'Active',
      count: activeEmotes.length,
      emptyMessage: 'No active emotes',
      content: <EmoteGrid emotes={activeEmotes} />,
    },
    {
      id: 'added',
      label: 'Added',
      count: countEmotes(addedEmotes),
      emptyMessage: 'No added emotes',
      content: <EmoteDateGroups emotes={addedEmotes} />,
    },
    {
      id: 'removed',
      label: 'Removed',
      count: countEmotes(removedEmotes),
      emptyMessage: 'No removed emotes',
      content: <EmoteDateGroups emotes={removedEmotes} />,
    },
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
      {tabs.map(({ id, count, emptyMessage, content }) => (
        <TabsContent key={id} value={id}>
          <EmotesList count={count} emptyMessage={emptyMessage}>
            {content}
          </EmotesList>
        </TabsContent>
      ))}
    </Tabs>
  );
};

export default EmoteTabs;
