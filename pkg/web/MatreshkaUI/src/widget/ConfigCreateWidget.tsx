import {useState} from "react";
import {ConfigTypePrefix, CreateConfigRequest} from "@vervstack/matreshka";

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
    const [type, setType] = useState<ConfigTypePrefix>(ConfigTypePrefix.verv);

    const [func, setFunc] = useState<Promise<void> | undefined>();

    const [isCreating, setIsCreating] = useState(false);

    const {CreateConfig} = useApi();
    const {CloseDialog} = useDialog();

    function create() {
        const req = {
            configName: name,
            configType: type,
        } as CreateConfigRequest;
        setIsCreating(true)
        setFunc(CreateConfig(req)
            .finally(() => setIsCreating(false))
            .then(CloseDialog))
    }

    function clear() {
        setName("")
        setFunc(undefined)
    }

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
                        onChange={setName}
                    />

                    <SelectDropDown
                        label={'Type'}
                        options={[
                            ConfigTypePrefix.plain, ConfigTypePrefix.verv, ConfigTypePrefix.minio,
                            ConfigTypePrefix.pg, ConfigTypePrefix.nginx, ConfigTypePrefix.kv,
                        ]}
                        selectedOption={type}
                        onSelected={(s: string) => setType(s as ConfigTypePrefix)}
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
