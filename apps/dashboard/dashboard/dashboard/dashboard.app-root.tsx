import {BrowserRouter} from 'react-router-dom';
import {createRoot,} from 'react-dom/client';
import {Dashboard} from "./dashboard.js";

/**
 * comment this in for server-side rendering (ssr) and comment
 * out of the root.render() invocation below.
 */
// hydrateRoot(
//   document.getElementById("root") as HTMLElement,
//   <BrowserRouter>
//     <Dashboard />
//   </BrowserRouter>
// );

if (import.meta.hot) {
    import.meta.hot.accept();
}

/**
 * mounting for client side rendering.
 */
const container = document.getElementById('root');
const root = createRoot(container!);

root.render(
    <BrowserRouter>
        <Dashboard/>
    </BrowserRouter>
);