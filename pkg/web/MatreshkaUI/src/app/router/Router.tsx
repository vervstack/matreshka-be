import {createBrowserRouter} from "react-router-dom";
import MainLayout from "@/app/router/MainLayout.tsx";
import ErrorPage from "@/pages/ErrorPage.tsx";
import HomePage from "@/pages/HomePage.tsx";

export enum Routes {
    Home = '/',
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
        ]
    }
])


export default router;
