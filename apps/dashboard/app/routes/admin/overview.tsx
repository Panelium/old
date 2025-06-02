import React from 'react';
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '~/components/ui/card';

export default function AdminOverviewPage() {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Admin Overview</CardTitle>
                <CardDescription>System-wide overview and statistics.</CardDescription>
            </CardHeader>
            <CardContent>
                <p>Manage users, servers, nodes, and system settings from here.</p>
                {/* Placeholder for system stats, quick links, etc. */}
            </CardContent>
        </Card>
    );
}

