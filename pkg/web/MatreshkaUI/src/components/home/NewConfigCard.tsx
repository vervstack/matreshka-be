import cn from "classnames";

import PencilIcon from "@/assets/icons/Pencil.svg";

import cls from "@/components/home/Card.module.css";
import {DialogType, useDialog} from "@/app/hooks/dialog/Dialog.tsx";

export default function NewConfigCard() {
    const {OpenDialog} = useDialog()

    return (
        <div
            className={cn(cls.CardContainer, cls.NewConfigCard)}
            onClick={() => OpenDialog(DialogType.CreateConfig)}
        >
            New
            <img
                className={cls.NewIcon}
                src={PencilIcon} alt={''}/>
        </div>
    )
}
