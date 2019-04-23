// Libraries
import React, {PureComponent} from 'react'

// Components
import SelectableBucket from 'src/authorizations/components/SelectableBucket'

interface Props {
  buckets: {id: string; name: string}[]
  onSelect: (id: string) => void
}

class SelectableBuckets extends PureComponent<Props> {
  render() {
    const {buckets, onSelect} = this.props
    return (
      <div>
        {buckets.map(b => {
          return <SelectableBucket bucket={b} onSelect={onSelect} key={b.id} />
        })}
      </div>
    )
  }
}

export default SelectableBuckets
