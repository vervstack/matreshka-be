import {useEffect, useState} from "react";

import cls from "@/segments/Loader.module.css";
import RetryIcon from "@/assets/icons/Retry.svg";

import {useToaster} from "@/app/hooks/toaster/Toaster.ts";
import ActionButton from "@/components/shared/ActionButton.tsx";


interface LoaderWrapperProps {
    children: React.JSX.Element[] | React.JSX.Element;

    load?: Promise<void>;

    onClickTryAgain?: () => void;
}

export default function LoaderWrapper({load, children, onClickTryAgain}: LoaderWrapperProps) {
    const [isLoading, setIsLoading] = useState(false);
    const [err, setErr] = useState<Error | undefined>();
    const toaster = useToaster()

    useEffect(() => {
        if (err) {
            toaster.catchGrpc(err)
        }
    }, [err]);

    function onClickTryAgainLocal() {
        if (!onClickTryAgain) return

        onClickTryAgain()
        setErr(undefined)
        setIsLoading(false)
    }

    useEffect(() => {
        // TODO remove from eslint
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setErr(undefined);

        if (!load) {
            setIsLoading(false);
            return
        }

        setIsLoading(true);
        load
            .catch(setErr)
            .finally(() => setIsLoading(false));
    }, [load]);

    return (
        <div className={cls.LoaderContainer}>
            {isLoading && <div className={cls.Spinner}/>}
            {err ?
                <div className={cls.ErrorMessageBox}>
                    <div>{err.message}</div>
                    {onClickTryAgain && <ActionButton

                        onClick={onClickTryAgainLocal}
                        iconPath={RetryIcon}/>}
                </div>

             : children}


        </div>
    )
}
