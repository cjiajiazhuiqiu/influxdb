import React, {PureComponent, ReactElement} from 'react'
import FuncArg from 'src/ifql/components/FuncArg'
import {OnChangeArg} from 'src/types/ifql'
import {ErrorHandling} from 'src/shared/decorators/errors'
import {Func} from 'src/types/ifql'
import {funcNames} from 'src/ifql/constants'
import Join from 'src/ifql/components/Join'

interface Props {
  func: Func
  bodyID: string
  onChangeArg: OnChangeArg
  declarationID: string
  onGenerateScript: () => void
  onDeleteFunc: () => void
  declarationsFromBody: string[]
}

@ErrorHandling
export default class FuncArgs extends PureComponent<Props> {
  public render() {
    const {
      func,
      bodyID,
      onChangeArg,
      onDeleteFunc,
      declarationID,
      onGenerateScript,
      declarationsFromBody,
    } = this.props
    const {name: funcName, id: funcID} = func
    if (funcName === funcNames.JOIN) {
      return (
        <div className="func-node--tooltip">
          <Join
            func={func}
            bodyID={bodyID}
            declarationID={declarationID}
            onChangeArg={onChangeArg}
            declarationsFromBody={declarationsFromBody}
            onGenerateScript={onGenerateScript}
          />
          <div className="func-node--buttons">
            <div
              className="btn btn-sm btn-danger func-node--delete"
              onClick={onDeleteFunc}
            >
              Delete
            </div>
            {this.build}
          </div>
        </div>
      )
    }
    return (
      <div className="func-node--tooltip">
        {func.args.map(({key, value, type}) => {
          return (
            <FuncArg
              key={key}
              type={type}
              argKey={key}
              value={value}
              bodyID={bodyID}
              funcID={funcID}
              funcName={funcName}
              onChangeArg={onChangeArg}
              declarationID={declarationID}
              onGenerateScript={onGenerateScript}
            />
          )
        })}
        <div className="func-node--buttons">
          <div
            className="btn btn-sm btn-danger func-node--delete"
            onClick={onDeleteFunc}
          >
            Delete
          </div>
          {this.build}
        </div>
      </div>
    )
  }

  get build(): ReactElement<HTMLDivElement> {
    const {func, onGenerateScript} = this.props
    if (func.name === funcNames.FILTER) {
      return (
        <div
          className="btn btn-sm btn-primary func-node--build"
          onClick={onGenerateScript}
        >
          Build
        </div>
      )
    }

    return null
  }
}
