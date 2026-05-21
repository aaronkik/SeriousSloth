import { useState } from 'react';
import { Badge } from '~/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs';
import type { ActiveEmote, RemovedEmote } from '~/lib/api/emotes-service';
import EmotesList from './emotes-list';

type Tab = 'active' | 'removed';

type Props = {
  activeEmotes: ActiveEmote[];
  removedEmotes: RemovedEmote[];
};

const EmoteTabs = ({ activeEmotes, removedEmotes }: Props) => {
  const [tab, setTab] = useState<Tab>('active');

  const tabs: Array<{ id: Tab; label: string; count: number; emotes: (ActiveEmote | RemovedEmote)[] }> = [
    { id: 'active', label: 'Active', count: activeEmotes.length, emotes: activeEmotes },
    { id: 'removed', label: 'Removed', count: removedEmotes.length, emotes: removedEmotes },
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
      {tabs.map(({ id, emotes }) => (
        <TabsContent key={id} value={id}>
          <EmotesList emotes={emotes} />
        </TabsContent>
      ))}
    </Tabs>
  );
};

export default EmoteTabs;
