import React, {PureComponent, ChangeEvent} from 'react'
import {connect} from 'react-redux'

// Components
import {Overlay, Input, Form} from 'src/clockface'
import {
  IconFont,
  ComponentColor,
  ComponentSpacer,
  AlignItems,
  FlexDirection,
  ComponentSize,
  Button,
  ButtonType,
} from '@influxdata/clockface'
import {withRouter, WithRouterProps} from 'react-router'

// Actions
import {createAuthorization} from 'src/authorizations/actions'

// Decorators
import {ErrorHandling} from 'src/shared/decorators/errors'
import {Authorization} from '@influxdata/influx'
import SelectableBuckets from './SelectableBuckets'
import GetResources, {ResourceTypes} from 'src/shared/components/GetResources'

interface DispatchProps {
  onCreateAuthorization: typeof createAuthorization
}

interface State {
  description: string
}

type Props = WithRouterProps & DispatchProps

@ErrorHandling
class BucketsTokenOverlay extends PureComponent<Props, State> {
  public state = {description: ''}

  render() {
    const {description} = this.state

    return (
      <Overlay visible={true}>
        <Overlay.Container>
          <Overlay.Heading
            title="Generate Read/Write Token"
            onDismiss={this.handleDismiss}
          />
          <Overlay.Body>
            <Form onSubmit={this.handleSave}>
              <ComponentSpacer
                alignItems={AlignItems.Center}
                direction={FlexDirection.Column}
                margin={ComponentSize.Large}
              >
                <Form.Element label="Description">
                  <Input
                    placeholder="Describe this new token"
                    value={description}
                    onChange={this.handleInputChange}
                  />
                </Form.Element>
                <Form.Element label="">
                  <GetResources resource={ResourceTypes.Buckets}>
                    <SelectableBuckets />
                  </GetResources>
                </Form.Element>
                <ComponentSpacer
                  alignItems={AlignItems.Center}
                  direction={FlexDirection.Row}
                  margin={ComponentSize.Small}
                >
                  <Button
                    text="Cancel"
                    icon={IconFont.Remove}
                    onClick={this.handleDismiss}
                  />

                  <Button
                    text="Save"
                    icon={IconFont.Checkmark}
                    color={ComponentColor.Success}
                    type={ButtonType.Submit}
                  />
                </ComponentSpacer>
              </ComponentSpacer>
            </Form>
          </Overlay.Body>
        </Overlay.Container>
      </Overlay>
    )
  }

  private handleSave = async () => {
    const {
      params: {orgID},
      onCreateAuthorization,
    } = this.props

    const token: Authorization = {
      orgID,
      description: this.state.description,
      permissions: [],
    }

    // await onCreateAuthorization(token)

    this.handleDismiss()
  }

  private handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const {value} = e.target

    this.setState({description: value})
  }

  private handleDismiss = () => {
    const {
      router,
      params: {orgID},
    } = this.props

    router.push(`/orgs/${orgID}/tokens`)
  }
}

const mdtp: DispatchProps = {
  onCreateAuthorization: createAuthorization,
}

export default connect<{}, DispatchProps, {}>(
  null,
  mdtp
)(withRouter(BucketsTokenOverlay))
