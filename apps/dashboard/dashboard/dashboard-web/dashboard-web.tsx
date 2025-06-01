import {Route, Routes} from 'react-router-dom';
import {Link, NavigationProvider} from '@bitdesign/sparks.navigation.link';
import {Header} from '@bitdesign/sparks.layout.header';
import {Logo} from '@bitdesign/sparks.content.logo';
import {AcmeTheme} from '@acme/design.acme-theme';
import {Announcements} from '@panelium/dashboard.ui.announcements';

export function DashboardWeb() {
    return (
        <AcmeTheme>
            <NavigationProvider>
                <Header logo={<Logo src='https://static.bit.dev/extensions-icons/acme.svg' name='Acme' slogan='Inc.'/>}>
                    <Link href='/'>Investors</Link>
                    <Link href='/'>Onboarding</Link>
                </Header>
                <Routes>
                    <Route path="/" element={<Announcements/>}/>
                </Routes>
            </NavigationProvider>
        </AcmeTheme>
    );
}
