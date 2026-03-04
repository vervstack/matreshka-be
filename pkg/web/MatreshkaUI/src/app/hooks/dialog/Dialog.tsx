import {create} from "zustand";

export enum DialogType {
    CreateConfig = 'CreateConfig',
}

export interface DialogManager {
    CurrentDialog: DialogType | null;

    OpenDialog(d: DialogType): void;

    CloseDialog(): void;
}


export const useDialog =
    create<DialogManager>((set, _) => ({
        CurrentDialog: null,

        OpenDialog(d: DialogType) {
            set({CurrentDialog: d});
        },

        CloseDialog() {
            set({CurrentDialog: null})
        }
    }))
