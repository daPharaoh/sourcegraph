import * as React from 'react'
import { FilterTypes, filterTypeKeys } from '../../../../../shared/src/search/interactive/util'
import { defaultFilterTypes } from './AddFilterRow'
import { FilterTypesToProseNames } from './filters'

interface Props {
    onAddNewFilter: (filterType: FilterTypes) => void
}

export class AddFilterDropdown extends React.PureComponent<Props> {
    private onAddNewFilter = (e: React.ChangeEvent<HTMLSelectElement>): void => {
        this.props.onAddNewFilter(e.target.value as FilterTypes)
    }

    public render(): JSX.Element | null {
        return (
            <select
                className="form-control add-filter-dropdown e2e-filter-dropdown"
                onChange={this.onAddNewFilter}
                value="default"
            >
                <option value="default" disabled={true}>
                    Add filter…
                </option>
                {filterTypeKeys
                    .filter(filter => !defaultFilterTypes.includes(filter) && filter !== FilterTypes.case)
                    .map(filter => (
                        <option key={filter} value={filter} className={`e2e-filter-dropdown-option-${filter}`}>
                            {FilterTypesToProseNames[filter]}
                        </option>
                    ))}
            </select>
        )
    }
}
