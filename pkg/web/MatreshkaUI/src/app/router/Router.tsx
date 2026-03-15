import {createBrowserRouter} from "react-router-dom";
import MainLayout from "@/app/router/MainLayout.tsx";
import ErrorPage from "@/pages/ErrorPage.tsx";
import HomePage from "@/pages/HomePage.tsx";
import ConfigPage from "@/pages/ConfigPage.tsx";

export enum Routes {
    Home = '/',
    Config = '/:configName',
}

const router = createBrowserRouter([
    {
        path: Routes.Home,
        element: <MainLayout/>,
        errorElement: <ErrorPage/>,
        children: [
            {
                index: true,
                element: <HomePage/>,
            },
            {
                path: Routes.Config,
                element: <ConfigPage/>,
            },
        ]
    }
])


export default router;
