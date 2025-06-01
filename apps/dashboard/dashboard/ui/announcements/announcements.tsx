import type {ReactNode} from 'react';
import classNames from 'classnames';
import {Heading} from '@bitdesign/sparks.typography.heading';
import {Text} from '@bitdesign/sparks.typography.text';
import {useAnnouncements} from '@panelium/dashboard.hooks.use-announcements';
import {Flex} from '@bitdesign/sparks.layout.flex';
import {Card} from '@acme/design.content.card';
import styles from './announcements.module.scss';

export type AnnouncementsProps = {
    /**
     * children for the announcements
     */
    children?: ReactNode;

    /**
     * announcements to use
     */
    className?: string;
};

export function Announcements({className}: AnnouncementsProps) {
    const announcements = useAnnouncements();
    if (!announcements) return null;

    return (
        <Flex className={classNames(styles.announcement, className)}>
            {announcements.map((announcement, key) => {
                return (
                    <Card key={key}>
                        <Heading level={2}>{announcement.title}</Heading>
                        <Text>{announcement.date.toLocaleString()}</Text>
                    </Card>
                );
            })}
        </Flex>
    );
}
