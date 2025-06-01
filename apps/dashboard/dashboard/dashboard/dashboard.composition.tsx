import {MemoryRouter} from 'react-router-dom';
import {Dashboard} from "./dashboard.js";

export const DashboardBasic = () => {
    return (
        <MemoryRouter>
            <Dashboard/>
        </MemoryRouter>
    );
}