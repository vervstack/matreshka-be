import {Outlet} from "react-router-dom";

import Toaster from "@/segments/Toaster.tsx";
import Settings from "@/segments/Settings.tsx";

import Dialog from "@/segments/Dialog.tsx";

export default function MainLayout() {
    return (
        <div
            style={{
                display: "flex",
                flexDirection: "column-reverse",
            }}
        >
            <Outlet/>
            <Settings/>

            <Dialog/>
            <Toaster/>
        </div>
    )
}
