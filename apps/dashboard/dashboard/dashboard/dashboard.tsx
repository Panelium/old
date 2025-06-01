import {Route, Routes} from 'react-router-dom';

export function Dashboard() {
    return (
        <Routes>
            <Route path="/" element={<h1>Hello World!</h1>}/>
        </Routes>
    );
}