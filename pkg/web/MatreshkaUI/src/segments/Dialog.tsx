import cls from "@/segments/Dialog.module.css";

import ConfigCreateWidget from "@/widget/ConfigCreateWidget.tsx";
import {useDialog, DialogType} from "@/app/hooks/dialog/Dialog.tsx";

export default function Dialog() {
    const {CurrentDialog, CloseDialog} = useDialog();

    let dialogElement = null;

    if (CurrentDialog == DialogType.CreateConfig) {
        dialogElement = (
            <div className={cls.ConfigCreateDialogContainer}>
                <ConfigCreateWidget/>
            </div>
        )
    }
    if (!dialogElement) return null;

    return (
        <div
            className={cls.DialogContainer}
            onMouseDown={(e) => {
                if (e.target === e.currentTarget) CloseDialog();
            }}
        >
            <div
                onMouseDown={(e) => {
                    e.stopPropagation();
                }}
            >
                {dialogElement}
            </div>
        </div>
    );
}
