import React, {useEffect} from 'react';
import {Outlet, useNavigate} from 'react-router-dom';
import {useSession} from '~/providers/SessionProvider';

export default function AuthLayout() {
    const {authenticated} = useSession();
    const navigate = useNavigate();

    useEffect(() => {
        if (authenticated) {
            navigate('/', {replace: true});
        }
    }, [authenticated, navigate]);

    return (
        <div className="flex min-h-screen flex-col items-center justify-center bg-chart-background">
            <Outlet/>
        </div>
    );
}
