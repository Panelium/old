import {MockProvider} from '@acme/acme.testing.mock-provider';
import {DashboardWeb} from "./dashboard-web.js";

export const DashboardWebBasic = () => {
    return (
        <MockProvider noTheme>
            <DashboardWeb/>
        </MockProvider>
    );
}
