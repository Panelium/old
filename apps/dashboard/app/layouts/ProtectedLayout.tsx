import {useLocation, useNavigate} from "react-router-dom";
import {useEffect} from "react";
import {useSession} from "~/providers/SessionProvider";
import {Outlet} from "react-router";

export default function ProtectedLayout() {
    const {authenticated} = useSession();
    const location = useLocation();
    const navigate = useNavigate();

    const showPage = authenticated || location?.pathname?.startsWith("/auth");

    useEffect(() => {
        if (!showPage) {
            navigate("/auth", {replace: true});
            return;
        }
    }, [authenticated, location.pathname, navigate]);

    return (
        <>
            {showPage && <Outlet/>}
        </>
    );
}