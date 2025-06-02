import React from 'react';
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '~/components/ui/card';

export default function AdminModulesPage() {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Module Management</CardTitle>
                <CardDescription>Install, configure, and manage system modules.</CardDescription>
            </CardHeader>
            <CardContent>
                <p>This is where you will manage modules for the application.</p>
                {/* Placeholder for module listing, installation, etc. */}
            </CardContent>
        </Card>
    );
}

