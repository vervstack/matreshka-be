import {create} from "zustand";

export enum DialogType {
    CreateConfig = 'CreateConfig',
}

export interface DialogManager {
    CurrentDialog: DialogType | null;

    IsClickOffClosesDialog: boolean;

    LockClosing(): void;

    UnlockClosing(): void;

    OpenDialog(d: DialogType): void;

    CloseDialog(): void;
}


export const useDialog =
    create<DialogManager>((set, get) => ({
        CurrentDialog: null,
        IsClickOffClosesDialog: true,

        LockClosing() {
            set({IsClickOffClosesDialog: false})
        },
        UnlockClosing() {
            set({IsClickOffClosesDialog: true})
        },

        OpenDialog(d: DialogType) {
            set({CurrentDialog: d});
        },

        CloseDialog() {
            const {IsClickOffClosesDialog} = get()
            if (IsClickOffClosesDialog) {
                set({CurrentDialog: null})
            }
        }
    }))
