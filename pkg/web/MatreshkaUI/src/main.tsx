import {createRoot} from 'react-dom/client'
import {Tooltip} from "react-tooltip";
import {RouterProvider} from "react-router-dom";

import '@/index.module.css'

import router from "@/app/router/Router.tsx";

createRoot(document.getElementById('root')!)
    .render(
        <div>
            <link href="@/assets/font/Comfortaa.ttf" rel="stylesheet"/>

            <RouterProvider router={router}/>
            <Tooltip id={"root-tooltip"}/>
        </div>,
    )
