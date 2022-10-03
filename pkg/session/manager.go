package session

import "github.com/igevin/sepweb/pkg/context"

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

func (m *Manager) GetSession(ctx *context.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any)
	}
	val, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return val.(Session), nil
	}
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.SessCtxKey] = sess
	return sess, nil
}

func (m *Manager) InitSession(ctx *context.Context, id string) (Session, error) {
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *Manager) RefreshSession(ctx *context.Context) (Session, error) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	if err = m.Refresh(ctx.Req.Context(), sess.ID()); err != nil {
		return nil, err
	}
	if err = m.Inject(sess.ID(), ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *Manager) RemoveSession(ctx *context.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	if err = m.Store.Remove(ctx.Req.Context(), sess.ID()); err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
