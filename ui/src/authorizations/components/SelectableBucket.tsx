// Libraries
import React, {PureComponent} from 'react'

interface Props {
  bucket: {id: string; name: string}
  onSelect: (id: string) => void
}

class SelectableBucket extends PureComponent<Props> {
  render() {
    const {bucket} = this.props

    return <div onClick={this.handleClick}>{bucket}</div>
  }

  private handleClick = () => {
    return this.props.bucket.id
  }
}

export default SelectableBucket
