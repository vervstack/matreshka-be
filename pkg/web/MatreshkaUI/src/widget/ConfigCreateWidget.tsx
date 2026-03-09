import {useEffect, useState} from "react";
import {ConfigType, CreateConfigRequest} from "@vervstack/matreshka";

import cls from "@/widget/ConfigCreateWidget.module.css";
import ConfigFile from "@/assets/icons/ConfigFile.svg";

import Input from "@/components/shared/Input.tsx";
import ActionButton from "@/components/shared/ActionButton.tsx";

import SelectDropDown from "@/components/shared/SelectDropDown.tsx";
import {useApi} from "@/app/hooks/api/api.ts";
import LoaderWrapper from "@/segments/LoaderWrapper.tsx";
import {useDialog} from "@/app/hooks/dialog/Dialog.tsx";


export default function ConfigCreateWidget() {
    const [name, setName] = useState('');
    const [type, setType] = useState<ConfigType>(ConfigType.verv);

    const [func, setFunc] = useState<Promise<void> | undefined>();

    const [isCreating, setIsCreating] = useState(false);

    const {CreateConfig} = useApi();
    const {CloseDialog} = useDialog();

    function create() {
        if (isCreating) return

        const req = {
            configName: name,
            configType: type,
        } as CreateConfigRequest;

        setIsCreating(true)

        setFunc(
            CreateConfig(req)
            .finally(() => setIsCreating(false))
            .then(CloseDialog))
    }

    function clear() {
        setName("")
        setFunc(undefined)
    }

    function handleNameChange(s: string) {
        s = s.replace(" ", "_")
        setName(s)
    }

    useEffect(() => {
        const handleKeyPress = (event: KeyboardEvent) => {
            if (event.key === 'Enter') {
                create();
            }
        };

        window.addEventListener('keydown', handleKeyPress);
        return () => {
            window.removeEventListener('keydown', handleKeyPress);
        };
    }, []);

    return (
        <div className={cls.ConfigCreateWidgetContainer}>
            <LoaderWrapper
                load={func}
                onClickTryAgain={clear}
            >
                <div className={cls.Content}>
                    <div className={cls.Header}>Create new config</div>
                    <Input
                        label={'Name'}
                        inputValue={name}
                        onChange={handleNameChange}
                    />

                    <SelectDropDown
                        label={'Type'}
                        options={[
                            ConfigType.plain, ConfigType.verv, ConfigType.minio,
                            ConfigType.pg, ConfigType.nginx, ConfigType.kv,
                        ]}
                        selectedOption={type}
                        onSelected={(s: string) => setType(s as ConfigType)}
                    />

                    <div className={cls.ButtonsWrapper}>
                        <ActionButton
                            tooltip={'Create new config'}
                            iconPath={ConfigFile}
                            alt={'Create'}
                            onClick={isCreating ? undefined : create}
                        />
                    </div>
                </div>
            </LoaderWrapper>
        </div>
    )
}
