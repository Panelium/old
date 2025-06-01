import ReactDOMServer from "react-dom/server";
import {StaticRouter} from "react-router-dom/server.js";
import {DashboardWeb} from "./dashboard-web.js";

interface IRenderProps {
    path: string;
}

export const render = async ({path}: IRenderProps) => {
    return ReactDOMServer.renderToString(
        <StaticRouter location={path}>
            <DashboardWeb/>
        </StaticRouter>
    );
};

/**
 * implement loadScripts() to inject scripts to the head
 * during SSR.
 */
// export const loadScripts = async () => {
//   return '<script></script>';
// }
