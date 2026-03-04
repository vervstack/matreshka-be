import cls from "@/segments/Dialog.module.css";

import ConfigCreateWidget from "@/widget/ConfigCreateWidget.tsx";
import {useDialog, DialogType} from "@/app/hooks/dialog/Dialog.tsx";
import {useRef} from "react";

export default function Dialog() {
    const dialog = useDialog()

    let dialogElement = null;

    const dialogRef = useRef(null);

    if (dialog.CurrentDialog == DialogType.CreateConfig) {
        dialogElement = (
            <div className={cls.ConfigCreateDialogContainer}>
                <ConfigCreateWidget/>
            </div>
        )
    }
    if (!dialogElement) return

    return (
        <div className={cls.DialogContainer}>
            <div ref={dialogRef}>
                {dialogElement}
            </div>
        </div>
    );
}
