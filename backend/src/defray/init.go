package defray

import (
	errors "github.com/wuntsong-org/wterrors"
)

func InitDefray() errors.WTError {
	err := startAllNotify()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = startDistribution()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
