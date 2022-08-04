// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { useFormContext } from 'react-hook-form';
import { CdsSelect } from '@cds/react/select';
import * as yup from 'yup';
import { AWSVirtualMachine } from '../../../../swagger-api';

interface ImageProps {
    osImageTitle: string;
    field: string;
    onOsImageSelected: (osImage: string, fieldName?: string) => void;
    images: AWSVirtualMachine[];
}
export function osImagesValidation() {
    return yup.string().nullable().required('Please select an OS image');
}

function OsImageSelect(props: ImageProps) {
    const { osImageTitle, field, images, onOsImageSelected } = props;
    const handleOsImageSelect = (event: ChangeEvent<HTMLSelectElement>) => {
        onOsImageSelected(event.target.value || '', field);
    };
    const {
        register,
        formState: { errors },
    } = useFormContext();

    const fieldError = errors[field];
    return (
        <div cds-layout="m:lg">
            <h1>{osImageTitle}</h1>
            <CdsSelect layout="compact" controlWidth="shrink">
                <label>OS Image with Kubernetes </label>
                <select {...register(field)} onChange={handleOsImageSelect}>
                    {images.map((image) => (
                        <option key={image.name}>{image.name}</option>
                    ))}
                </select>
            </CdsSelect>
            {fieldError && (
                <div>
                    &nbsp;<CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>
                </div>
            )}
        </div>
    );
}

export default OsImageSelect;
