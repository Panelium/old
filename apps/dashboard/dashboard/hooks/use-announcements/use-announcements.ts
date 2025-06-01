import {gql, useQuery} from '@apollo/client';
import {useIsMock} from '@acme/acme.testing.mock-provider';
import {Announcement, mockAnnouncements} from '@panelium/dashboard.entities.announcement';

// define a graphql query
export const LIST_ANNOUNCEMENTS = gql`
  query LIST_ANNOUNCEMENTS {
    listAnnouncements {
      title
      date
    }
  }
`;

/**
 * fetch list of announcements.
 */
export function useAnnouncements(): undefined | null | Announcement[] {
    const results = useQuery(LIST_ANNOUNCEMENTS);
    const isMock = useIsMock();
    if (isMock) return mockAnnouncements();

    if (!results.data || results.loading) return undefined;
    if (!results?.data?.listAnnouncements) return null;

    return results.data.listAnnouncements.map((announcement) => {
        return Announcement.from(announcement);
    });
}
