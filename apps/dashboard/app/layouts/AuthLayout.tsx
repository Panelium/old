import React from 'react';
import {Outlet} from 'react-router-dom'; // Changed from @remix-run/react

export default function AuthLayout() {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center bg-gray-100 dark:bg-gray-900">
            <Outlet/>
        </div>
    );
}
