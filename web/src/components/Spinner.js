import { memo } from 'react';
import { CgSpinner } from 'react-icons/cg';

const Spinner = () => {
    return (
        <CgSpinner className="animate-spin inline-block align-middle loader ease-linear rounded-full h-6 w-6" />
    );
};

export default memo(Spinner);
